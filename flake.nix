{
  description = "Flake to develop arus locally";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
    process-compose-flake.url = "github:Platonic-Systems/process-compose-flake";
    services-flake.url = "github:juspay/services-flake";
  };

  outputs = { flake-utils, nixpkgs, ... } @ input:
    flake-utils.lib.eachDefaultSystem (system:
      let
        goVersion = 24;

        overlays = [
          (self: super: {
            go = super."go_1_${toString goVersion}";
          })
        ];

        pkgs = import nixpkgs { inherit system overlays; config.allowUnfree = true; config.allowBroken = true; };

        service = import ./compose.nix ({ inherit pkgs; } // input);

      in
      {
        packages = {
          live = service.mkProcess service.compose.live;
          runServices = service.mkProcess service.compose.infra;
        };

        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            just
            postgresql_jit
            sqlc
            go
            air
            nodejs
            nodePackages.pnpm
            atlas
          ];
        };
      }
    );
}
