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

              # Required until golangci-lint is built with Go 1.19
              (
                final: prev: {
                  golangci-lint = prev.golangci-lint.override {
                    buildGoModule = final.buildGo119Module;
                  };
                }
              )
            ];
          };

          buildDeps = with pkgs; [ git go_1_19 gnumake ];
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
