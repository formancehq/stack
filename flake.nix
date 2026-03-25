{
  description = "Formance Stack - OpenAPI spec build and publish environment";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2511";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forEachSupportedSystem = f:
        nixpkgs.lib.genAttrs supportedSystems (system:
          let
            pkgs = import nixpkgs { inherit system; };
          in
          f { pkgs = pkgs; system = system; }
        );

      speakeasyVersion = "1.759.2";
      speakeasyPlatforms = {
        "x86_64-linux"   = "linux_amd64";
        "aarch64-linux"  = "linux_arm64";
        "x86_64-darwin"  = "darwin_amd64";
        "aarch64-darwin" = "darwin_arm64";
      };
      speakeasyHashes = {
        "x86_64-linux"   = "9234e2e9138f03ac18f0ca034d0c5a0a7b8749ea91b16814b4a643afd680d8fd";
        "aarch64-linux"  = "ba92a8c5799ed14acba94d317694ed32e35883e9439a07b28c58f7c8c0ea16f5";
        "x86_64-darwin"  = "b4cfe13627e8822718b5820c68898f51b6381e604c9578650c9b0c3f40f800b0";
        "aarch64-darwin" = "dda057dbbd83bdaa47f9ccf3311e455013d957d11f45d8336b97b91ba2a56d6d";
      };

    in
    {
      packages = forEachSupportedSystem ({ pkgs, system }:
        {
          speakeasy = pkgs.stdenv.mkDerivation {
            pname = "speakeasy";
            version = speakeasyVersion;

            src = pkgs.fetchurl {
              url = "https://github.com/speakeasy-api/speakeasy/releases/download/v${speakeasyVersion}/speakeasy_${speakeasyPlatforms.${system}}.zip";
              sha256 = speakeasyHashes.${system};
            };

            nativeBuildInputs = [ pkgs.unzip ];
            dontUnpack = true;

            installPhase = ''
              mkdir -p $out/bin
              unzip $src
              install -m755 speakeasy $out/bin/
            '';

            name = "speakeasy";
          };
        }
      );

      devShells = forEachSupportedSystem ({ pkgs, system }:
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.jq
              pkgs.just
              pkgs.nodejs_22
              pkgs.wget
              pkgs.yq-go
              self.packages.${system}.speakeasy
            ];
          };
        }
      );
    };
}
