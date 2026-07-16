{
  description = "Collection of small Go libraries used by nais K8s operators";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.treefmt-nix = {
    url = "github:numtide/treefmt-nix";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    inputs:
    inputs.flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import inputs.nixpkgs {
          localSystem = { inherit system; };
          overlays = [
            (
              final: prev:
              let
                version = "1.26.5";
                newerGoVersion = prev.go_latest.overrideAttrs (old: {
                  inherit version;
                  src = prev.fetchurl {
                    url = "https://go.dev/dl/go${version}.src.tar.gz";
                    hash = "sha256-SVvkvIcXasVnOS5bQRar2YRm0z17SdQedkzMaXay3EI=";
                  };
                });
                nixpkgsVersion = prev.go_latest.version;
                newVersionNotInNixpkgs = -1 == builtins.compareVersions nixpkgsVersion version;
              in
              {
                go_latest = if newVersionNotInNixpkgs then newerGoVersion else prev.go_latest;
                buildGoModule = prev.buildGoModule.override { go = final.go_latest; };
              }
            )
          ];
        };
        inherit (pkgs) lib;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            go-mockery
            gofumpt
            gopls
            gotools
            kustomize
          ];
        };
        formatter = inputs.treefmt-nix.lib.mkWrapper pkgs (
          {
            projectRootFile = "flake.nix";
            settings.global.excludes = [
              "*.md"
              ".gitattributes"
            ];
          }
          // {
            programs =
              lib.genAttrs
                [
                  # General
                  "shellcheck"
                  "dos2unix"

                  # go
                  "gofumpt"

                  # nix
                  "statix"
                  "nixfmt"
                  "deadnix"
                ]
                (_: {
                  enable = true;
                });
          }
        );
      }
    );
}
