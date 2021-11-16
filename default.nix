{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "05s96aglpsg32f6md6rsrda0r3pgg87b7fzj13a37zrs0bsirj73";
  src = ./.;
  doCheck = false;
}
