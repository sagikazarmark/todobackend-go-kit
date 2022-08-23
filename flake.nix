{
  description = "A simple Todo-Backend application written using Go kit";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    goflake.url = "github:sagikazarmark/go-flake";
    goflake.inputs.nixpkgs.follows = "nixpkgs";
    gobin.url = "github:sagikazarmark/go-bin-flake";
    gobin.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, goflake, gobin, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
        let
          pkgs = import nixpkgs {
            inherit system;

            overlays = [
              goflake.overlay

              (
                final: prev: {
                  golangci-lint = gobin.packages.${system}.golangci-lint-bin;
                  dagger = prev.buildGo118Module rec {
                    pname = "dagger";
                    version = "0.2.11";

                    src = prev.fetchFromGitHub {
                      owner = "dagger";
                      repo = "dagger";
                      rev = "v${version}";
                      sha256 = "sha256-jkH1OrddLUMSj0Hs5T9jyVVR9F5x7jzIZ8HYixA0x2s=";
                    };

                    vendorSha256 = "sha256-4GmdgyoqArvjJsQsVjwaxlvOMwYHUTiuD1jOzW8DPKQ=";

                    proxyModule = true;

                    subPackages = [
                      "cmd/dagger"
                    ];

                    ldflags = [ "-s" "-w" "-X go.dagger.io/dagger/version.Revision=${version}" ];
                  };
                }
              )
            ];
          };

          buildDeps = with pkgs; [ git go_1_18 gnumake ];
          devDeps = with pkgs; buildDeps ++ [
            golangci-lint
            gotestsum
            goreleaser
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
            protoc-gen-go-kit
            # gqlgen
            openapi-generator-cli
            dagger
          ];
        in
          { devShell = pkgs.mkShell { buildInputs = devDeps; }; }
    );
}
