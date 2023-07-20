# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles }:
let
  shaMap = {
    x86_64-linux = "e51f209dc332a9f5d54880b8bc0d2f09b9fc912e608ed3220021fcd65ba22b9a";
    aarch64-linux = "afc000e9541fc99ac6210fec0429a676b974115944c0461496a34bb4291d8b47";
    x86_64-darwin = "918177be33d3780536684082fffa5d45ffa70e9bdeea93d365482bf4295555b8";
    aarch64-darwin = "cbf6b05fdabff5fe663545894b9c098bff3d947eaba0a7286642e1d52806193e";
  };

  urlMap = {
    x86_64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.10.1/moon-x86_64-unknown-linux-gnu";
    aarch64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.10.1/moon-aarch64-unknown-linux-gnu";
    x86_64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.10.1/moon-x86_64-apple-darwin";
    aarch64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.10.1/moon-aarch64-apple-darwin";
  };
in pkgs.stdenv.mkDerivation {
  pname = "moon";
  version = "1.10.1";
  src = fetchurl {
    url = urlMap.${system};
    sha256 = shaMap.${system};
  };

  sourceRoot = ".";

  dontUnpack = true;

  nativeBuildInputs = [ installShellFiles ];

  installPhase = ''
    install -D $src $out/bin/moon
    chmod a+x $out/bin/moon
  '';

  system = system;

  meta = with lib; {
    description = "A task runner and repo management tool for the web ecosystem, written in Rust.";
    homepage = "https://moonrepo.dev/moon";
    license = licenses.mit;

    platforms = [
      "aarch64-darwin"
      "aarch64-linux"
      "x86_64-darwin"
      "x86_64-linux"
    ];
  };
}
