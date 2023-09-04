# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles, unzip }:
let
  shaMap = {
    x86_64-linux = "3453b764cb8f5f51b8bbf6f5d3846408157546f9370f62d71f1480160b102806";
    aarch64-linux = "446ac0310add7e27e0276678bdd03b667fae1e6d5b4c1bf34d1607e04a641aec";
    i686-linux = "f69e22967b87fa680e2f4d2a600c3865b6ae226f08fc6898514fbfa925a0fb8e";
    x86_64-darwin = "ef8a25c3b8048a02b0740537946f72268e2fab0ddb07872d5c39c659d8a240b2";
    aarch64-darwin = "54f39aef42b2504ff731940d224bcd94aaef9fb3c83a77bb9653ec49b00c0fb6";
  };

  urlMap = {
    i686-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.1/speakeasy_linux_386.zip";
    x86_64-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.1/speakeasy_linux_amd64.zip";
    aarch64-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.1/speakeasy_linux_arm64.zip";
    x86_64-darwin =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.1/speakeasy_darwin_amd64.zip";
    aarch64-darwin =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.1/speakeasy_darwin_arm64.zip";
  };
in pkgs.stdenv.mkDerivation {
  pname = "speakeasy";
  version = "1.78.1";
  src = fetchurl {
    url = urlMap.${system};
    sha256 = shaMap.${system};
  };

  sourceRoot = ".";

  nativeBuildInputs = [ installShellFiles unzip ];

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
