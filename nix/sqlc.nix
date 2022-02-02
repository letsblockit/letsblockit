{ buildGoModule, fetchFromGitHub, go_1_17 }:
buildGoModule.override { go = go_1_17; } {
  pname = "sqlc";
  version = "1.10.0";

  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v1.10.0";
    sha256 = "0fp56kxpaknffwjxpmaa5g7hmqdprc6adgw8bzw0a3c281r2lsdl";
  };

  vendorSha256 = "sha256-gxzmWjhGXACPLyOrquoCw6XN1vKqXDh7WrsYkxpHYkw=";
  runVend = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}
