{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "18l3immqs570pjzk3gwhjbjpf7i6s150qcrriz8gdz8wsnqq71hj";
  src = ./.;
  doCheck = false;
}
