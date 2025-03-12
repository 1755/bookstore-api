{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  dotenv.enable = true;

  # https://devenv.sh/packages/
  packages = [ 
    pkgs.git
    pkgs.docker
    pkgs.httpie
    pkgs.postgresql_17
   ];

  # https://devenv.sh/languages/
  languages.go = {
    enable = true;
  };

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  services = {
    postgres = {
      enable = true;
      package = pkgs.postgresql_17;
      initialDatabases = [
        {
          name = "postgres";
          user = "postgres";
          pass = "postgres";
        }
      ];
      listen_addresses = "127.0.0.1";
      port = 5432;
    };
  };

  # https://devenv.sh/scripts/
  scripts.hello.exec = ''
    echo hello from $GREET
  '';

  enterShell = ''
    unset PYTHONPATH;
    #go install github.com/nametake/golangci-lint-langserver@latest
    #go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0
    git --version
    go version
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  git-hooks.hooks = {
    gofmt.enable = true;
    golangci-lint.enable = true;
  };

  # See full reference at https://devenv.sh/reference/options/
}
