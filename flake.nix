{
  description = "Collection of small Go libraries used by nais K8s operators";

  # Flake inputs
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.treefmt-nix = {
    url = "github:numtide/treefmt-nix";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  # Flake outputs
  outputs =
    inputs:
    let
      goOverlay =
        final: prev:
        let
          nixpkgsGo = prev.go;
          firstSemverGeqSecond = first: second: builtins.compareVersions first second >= 0;
          goModFileVersion = builtins.elemAt (builtins.match ".*\ngo ([0-9a\.a-z]+).*" (builtins.readFile ./go.mod)) 0;

          projectGo =
            if (firstSemverGeqSecond nixpkgsGo.version goModFileVersion) then
              builtins.trace "Using nixpkgs' Go version: ${nixpkgsGo.version}" nixpkgsGo
            else
              (
                if (firstSemverGeqSecond upstream.version nixpkgsGo.version) then
                  builtins.warn "Using `upstream`'s Go version: ${upstream.version}, update nixpkgs?"
                    prev.go.overrideAttrs
                    (_: {
                      src = prev.fetchurl rec {
                        inherit (upstream) hash version;
                        url = "https://go.dev/dl/go${version}.src.tar.gz";
                      };
                    })
                else
                  builtins.abort "!!! UPDATE nixpkgs or `upstream` Go's referenced version to match the one `go.mod` requires!"
                # See below:
              );

          upstream.version = "1.23.0"; # Change us when above error is thrown
          upstream.hash = "sha256-Qreo6A2AXaoDAi7T/eQyHUw78smQoUQWXQHu7Nb2mcY="; # Change us when above error is thrown
        in
        {
          go = projectGo;
          buildGoModule = prev.buildGoModule.override { inherit (final) go; };
        };
      # Systems supported
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];

      # Helper to provide system-specific attributes
      forAllSystems =
        f:
        inputs.nixpkgs.lib.genAttrs allSystems (
          system:
          f rec {
            inherit (pkgs) lib;
            pkgs = import inputs.nixpkgs {
              inherit system;
              overlays = [ goOverlay ];
            };
          }
        );
    in
    {
      # Development environment output
      devShells = forAllSystems (
        { pkgs, ... }:
        {
          default = pkgs.mkShell {
            # The Nix packages provided in the environment
            packages = with pkgs; [
              go
              go-mockery
              gofumpt
              gopls
              gotools
            ];
          };
        }
      );
      formatter = forAllSystems (
        { pkgs, lib, ... }:
        inputs.treefmt-nix.lib.mkWrapper pkgs (
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
        )
      );
    };
}
