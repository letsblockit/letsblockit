{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  buildInputs = [ pkgs.sqlite ];
  buildFlags = "-tags libsqlite3";
  vendorSha256 = "1vxb38h7nwjvp0divx45kmmgbf4wljmsim9afvm3jvwg99466aai";
  src = ./.;
  doCheck = false;
}
