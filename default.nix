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
  vendorSha256 = "1j6skngagnddmhk980c7cmy30rkrlz1h9yjnk8rrnz8aday5gvx1";
  src = ./.;
  doCheck = false;
}
