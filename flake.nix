{
  inputs = {
    # self.submodules = true;
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default = pkgs.buildGoModule {
        name = "nuconv";
        system = builtins.currentSystem;
        meta.mainProgram = "nuconv";

        src = ./.;

        subPackages = [ "." ];
        vendorHash = "sha256-j3mvjbWefOhXGUk6euWBk+oPTujGVB9uDymdg7e/dy0=";
      };

      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/nuconv";
      };
    };
}
