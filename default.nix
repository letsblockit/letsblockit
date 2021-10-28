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
  vendorSha256 = "004k0ifb2k29xapd4zzhg7yl0dxz3jdg3xpkp2ssqx1kqnnzq5ik";
  src = ./.;
  doCheck = false;
}
