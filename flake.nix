{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    git-hooks.url = "github:cachix/git-hooks.nix";
    git-hooks.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = {
    self,
    nixpkgs,
    git-hooks,
  }: let
    forAllSystems = f:
      nixpkgs.lib.genAttrs ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"] (
        system: f nixpkgs.legacyPackages.${system}
      );
  in {
    checks = forAllSystems (pkgs: {
      pre-commit = git-hooks.lib.${pkgs.system}.run {
        src = ./.;
        hooks = {
          govet.enable = true;
          alejandra.enable = true;
          golangci-lint = {
            enable = true;
            name = "golangci-lint";
            entry = "${pkgs.golangci-lint}/bin/golangci-lint fmt";
            types = ["go"];
          };
        };
      };
    });

    devShells = forAllSystems (pkgs: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          golangci-lint
          docker

          (writeShellScriptBin "run" ''
            go run cmd/main.go
          '')

          (writeShellScriptBin "fmt" ''
            ${golangci-lint}/bin/golangci-lint fmt
          '')

          (writeShellScriptBin "dev" ''
            set -e

            if [ ! -f .env ]; then
              echo "err: .env file not found"
              exit 1
            fi

            mkdir -p "$(pwd)/data"

            docker build -t discord-quest-watcher:dev .

            docker run --rm \
              --name discord-quest-watcher-dev \
              --env-file .env \
              -v "$(pwd)/data:/data" \
              discord-quest-watcher:dev
          '')
        ];

        shellHook = self.checks.${pkgs.system}.pre-commit.shellHook;
      };
    });

    formatter = forAllSystems (pkgs: pkgs.alejandra);
  };
}
