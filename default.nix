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
  vendorSha256 = "1dbjnas15q2vcg5w7lrlz69ka1k1dygaf4glbww4zmf1y3ji68qm";
  src = ./.;
  doCheck = false;
}
