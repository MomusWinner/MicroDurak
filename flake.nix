{
  description = "A Nix-flake-based Go 1.24 development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1.*.tar.gz";
  inputs.defold.url = "github:spl3g/defold-flake";

  outputs = {
    self,
    defold,
    nixpkgs,
    ...
  }: let
    goVersion = 24; # Change this to update the whole stack

    supportedSystems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
    forEachSupportedSystem = f:
      nixpkgs.lib.genAttrs supportedSystems (system:
        f {
          pkgs = import nixpkgs {
            inherit system;
            overlays = [self.overlays.default];
          };
        });
  in {
    overlays.default = final: _: {
      go = final."go_1_${toString goVersion}";
    };

    devShells = forEachSupportedSystem ({pkgs}: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          # go (version is specified by overlay)
          go

          # goimports, godoc, lsp, etc.
          gotools
          gopls
          defold.packages."x86_64-linux".default

          # db stuff
          sqlc
          goose

          # generate rsa priv/pub keys for jwt
          openssl
        ];
      };
    });
  };
}
