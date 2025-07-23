{
  description = "A simple Go package";

  inputs.nixpkgs.url = "nixpkgs/nixos-25.05";

  outputs =
    { self, nixpkgs }:
    let
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.buildGoModule {
            inherit version;
            pname = "feedme";
            src = ./.;
            subPackages = [
              "cmd/..."
            ];
            vendorHash = null;
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              sqlc
              goose
            ];
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              go-tools
              gcc
            ];
          };
        }
      );
    };
}
