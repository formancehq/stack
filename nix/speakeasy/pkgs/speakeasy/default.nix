# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles }:
let
  shaMap = {
    x86_64-linux = "158a1c3303dee0307f12a04b9e17c97beef3081c4a58386b6b794c27cc7ae2d8";
    aarch64-linux = "4094df1bfc7b7b741cf85470185e65359d7339b65ef43351679812fcf831e4f3";
    i686-linux = "024835afffb1989662bed57facc2b1b5c23d299b93f4b6a90e5e9a745987b6aa";
    x86_64-darwin = "c4e2f99957f08c3cdf33b662798685f6c7c5db41a3e8d17265792d5e34624f56";
    aarch64-darwin = "fbdb38e8df1b86c943e319e7787c074f31d233af08bca12aaf9ae73b4cbaad5e";
  };

  urlMap = {
    i686-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.36.1/speakeasy_linux_386.tar.gz";
    x86_64-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.36.1/speakeasy_linux_amd64.tar.gz";
    aarch64-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.36.1/speakeasy_linux_arm64.tar.gz";
    x86_64-darwin =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.36.1/speakeasy_darwin_amd64.tar.gz";
    aarch64-darwin =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.36.1/speakeasy_darwin_arm64.tar.gz";
  };
in pkgs.stdenv.mkDerivation {
  pname = "speakeasy";
  version = "1.36.1";
  src = fetchurl {
    url = urlMap.${system};
    sha256 = shaMap.${system};
  };

  sourceRoot = ".";

  nativeBuildInputs = [ installShellFiles ];

  installPhase = ''
    mkdir -p $out/bin
    cp -vr ./speakeasy $out/bin/speakeasy
  '';

  system = system;

  meta = with lib; {
    description = "Speakeasy CLI makes validating API specs and generating idiomatic SDKs easy !";
    homepage = "https://speakeasyapi.dev";
    license = licenses.mit;

    platforms = [
      "aarch64-darwin"
      "aarch64-linux"
      "i686-linux"
      "x86_64-darwin"
      "x86_64-linux"
    ];
  };
}
