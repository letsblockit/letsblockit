{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "sqlc";
  version = "1.10.0";

  src = pkgs.fetchFromGitHub {
    owner = "kyleconroy";
    repo = "sqlc";
    rev = "v1.10.0";
    sha256 = "0fp56kxpaknffwjxpmaa5g7hmqdprc6adgw8bzw0a3c281r2lsdl";
  };
  vendorSha256 = "0k328wd9665vb9xkhp5aybbcv9f30bmamar35y7h0p2671dfc743";
  runVend = true; # pg_query_go ships the C headers in its module
  doCheck = false;
}
