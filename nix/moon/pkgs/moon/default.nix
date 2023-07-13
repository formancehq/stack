# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles }:
let
  shaMap = {
    x86_64-linux = "6a8f85070870258a1a875e4aba9d9fbf853172ee7bfe38ec9f8d88ae761ec281";
    aarch64-linux = "b86ddceed89615fe75f9c2c3e3d0002bf8fa11da784fa8887ac9dbb8324208a2";
    x86_64-darwin = "82558ca055bfee4fa44feb5e0c114153f11a3314ec8966e1d630bd36504fb597";
    aarch64-darwin = "9e4d2954a40675ebbb5789a46a8a3b50c1374575e45204ce4863ed2bc11fb4b8";
  };

  urlMap = {
    x86_64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.10.0/moon-x86_64-unknown-linux-gnu";
    aarch64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.10.0/moon-aarch64-unknown-linux-gnu";
    x86_64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.10.0/moon-x86_64-apple-darwin";
    aarch64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.10.0/moon-aarch64-apple-darwin";
  };
in pkgs.stdenv.mkDerivation {
  pname = "moon";
  version = "1.10.0";
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
