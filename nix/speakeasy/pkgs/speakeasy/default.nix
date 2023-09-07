# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles, unzip }:
let
  shaMap = {
    x86_64-linux = "c406a9520c02e1da00ff13368eb9a03bce4fae37207b9406b5b04cfea04f7733";
    aarch64-linux = "d167cd0856506b841656fae9ea4d1c788536abdb626af32c527569fcc62eb4eb";
    i686-linux = "f2ae0cbc6914c2eec119641d286618b2371200c0a7cf7b17a392a7fe511c66f3";
    x86_64-darwin = "8e3209bacfd758c0b96de204a0cd57b22e2b0cab63b297479eb855ef324a72bd";
    aarch64-darwin = "e302dcd6af5cacc41be461bad0f7ebc1b12bec8507f6169df18da7582b39d414";
  };

  urlMap = {
    i686-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.4/speakeasy_linux_386.zip";
    x86_64-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.4/speakeasy_linux_amd64.zip";
    aarch64-linux =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.4/speakeasy_linux_arm64.zip";
    x86_64-darwin =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.4/speakeasy_darwin_amd64.zip";
    aarch64-darwin =
      "https://github.com/speakeasy-api/speakeasy/releases/download/v1.78.4/speakeasy_darwin_arm64.zip";
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
    chmod +x $out/bin/speakeasy
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
