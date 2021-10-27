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
  vendorSha256 = "14p4q352sc68jb680jqqpdv64imba2q3k7bzbm3d8s4ljpk4gw4s";
  src = ./.;
  doCheck = false;
}
