# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles }:
let
  shaMap = {
    x86_64-linux = "c0c8fdefd9f5f53671308431c198973133c53e1543399e68855e5988f2d0566a";
    aarch64-linux = "47982a157d9b92d5c699cf078d2756d80d6703a4771b0cd8689be8f388551ad3";
    x86_64-darwin = "41a1c9f87b7843c8945556f0c73b0a4ae3b2bc8d4a9a1535a6c3f2330819cb62";
    aarch64-darwin = "94de430230cd3687b5444f2fbcfc548440432def7ce33c9c23dde0f1b71a2ce7";
  };

  urlMap = {
    x86_64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.6.1/moon-x86_64-unknown-linux-gnu";
    aarch64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.6.1/moon-aarch64-unknown-linux-gnu";
    x86_64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.6.1/moon-x86_64-apple-darwin";
    aarch64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.6.1/moon-aarch64-apple-darwin";
  };
in pkgs.stdenv.mkDerivation {
  pname = "moon";
  version = "1.6.1";
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
