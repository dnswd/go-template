{ pkgs, process-compose-flake, services-flake, ... }: rec {
  mkProcess = compose: (import process-compose-flake.lib { inherit pkgs; }).makeProcessCompose {
    name = "dev";
    modules = [
      services-flake.processComposeModules.default
      compose
    ];
  };

  mkPostgresInstance = { name, dataDir }: { config, ... }: {
    inherit dataDir;
    enable = true;
    package = pkgs.postgresql; # postgres 17
    extensions = exts: [
      exts.timescaledb
      exts.timescaledb_toolkit
      exts.pg_uuidv7
    ];
    settings = {
      listen_addresses = config.listen_addresses;
      port = config.port;
      unix_socket_directories = config.socketDir;
      hba_file = "${config.hbaConfFile}";
      shared_preload_libraries = "timescaledb";
    };
  };

  compose = rec {
    # This file contains recipe for spawning collection of processes akin to docker-compose.
    # However, instead of relying on docker, this recipe uses process-compose and nix
    # packages to guarantee reproducability. ach recipe can spawn multiple processes, and
    # each process can interact with each other as they are running without virt/container.
    #
    # References:
    # - https://github.com/Platonic-Systems/process-compose-flake/blob/main/nix/process-compose/default.nix
    # - https://github.com/mimsy-cms/mimsy/blob/b28ea78abde162c6c8276c91abcf68854f629311/flake-modules/process-compose.nix#L12
    # - https://github.com/graphprotocol/graph-node/blob/4de51ffc61cd010fee23858a00da892e3973bebb/flake.nix

    infra = { config, ... }: {
      settings.environment = {
        DB_URI = "${config.services.postgres."pgx".connectionURI { dbName = "postgres"; }}?sslmode=disable";
      };

      services.postgres."pgx" = mkPostgresInstance {
        name = "postgres";
        dataDir = "./.local/pg/data";
      };

      settings.processes."migration" = {
        command = ''
          just migrate
        '';
        depends_on = { "pgx".condition = "process_healthy"; };
      };

      settings.processes."export-infra-env" = {
        command = /* sh */ ''
          cat << EOF > ./.env.local
          ${builtins.concatStringsSep "\n" config.settings.environment}
          EOF

          echo "exported env to .env.local"
        '';
        depends_on = {
          # run when all service and init scripts done
          "pgx".condition = "process_healthy";
          "migration".condition = "process_completed_successfully";
        };
      };
    };

    live = { ... }: {
      imports = [ infra ];
      settings.processes."arus" = {
        command = ''
          just backend
        '';
        depends_on = {
          "export-infra-env".condition = "process_completed_successfully";
        };
      };
    };
  };
}

