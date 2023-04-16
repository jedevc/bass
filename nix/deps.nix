{ pkgs
}:
let
  markdownfmt = pkgs.buildGo118Module
    rec {
      name = "markdownfmt";

      vendorSha256 = "sha256-Ruj2Agh3P5Lt/KV3YcTpN54bbsbv/dQrLtKtIx08NE0=";

      src = pkgs.fetchFromGitHub {
        owner = "Kunde21";
        repo = "markdownfmt";
        rev = "f85609284a50d41a4c3a39cc112036929ba23af5";
        sha256 = "sha256-ggE6DK2iDAm8S5EZ8yRoYA5ISY21e3k9hKqOlZ517Ak=";
      };

      subPackages = [ "." ];
    };
in
with pkgs;
[
  # for running scripts
  bashInteractive
  # go building + testing
  go_1_20
  gcc
  gotestsum
  markdownfmt
  # lsp tests
  neovim
  # packing bass.*.(tgz|zip)
  gzip
  gnutar
  zip
  # git plumbing
  git
  # compressing shim binaries
  upx
  # for sanity checking that upx exists
  #
  # not needed by nix, but needed by Makefile
  which
  # for building in test image
  gnumake
  # for protoc
  protobuf
  protoc-gen-go
  protoc-gen-go-grpc
  # docs
  yarn
]
