# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{ system ? builtins.currentSystem, pkgs, lib, fetchurl, installShellFiles }:
let
  shaMap = {
    x86_64-linux = "de5e783840fedeed58d00601a00db029db8a63ff46277928bb04ff07b310eadd";
    aarch64-linux = "cecbd6e85f1baa7e79908710ea1a309e3db7cbbc7f1e3f69a417ee399de98a92";
    x86_64-darwin = "2abc5d1db5a4b4fa94f0e2773802a559acdfb5017e5ba2a3edfce538d9b542e8";
    aarch64-darwin = "c70fd2e4fb9c062173895b1da6105301771f190b54d28b152030457104ca555e";
  };

  urlMap = {
    x86_64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.8.3/moon-x86_64-unknown-linux-gnu";
    aarch64-linux =
      "https://github.com/moonrepo/moon/releases/download/v1.8.3/moon-aarch64-unknown-linux-gnu";
    x86_64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.8.3/moon-x86_64-apple-darwin";
    aarch64-darwin =
      "https://github.com/moonrepo/moon/releases/download/v1.8.3/moon-aarch64-apple-darwin";
  };
in pkgs.stdenv.mkDerivation {
  pname = "moon";
  version = "1.8.3";
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
