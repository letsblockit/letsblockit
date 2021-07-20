{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1nb7vj6sapzv585aiwlja00jzhr2z4k38n1al93lbbv9brik8ync";
  src = ./.;
}
