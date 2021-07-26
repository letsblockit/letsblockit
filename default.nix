{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "18d5wqy8r7xrrwp8v3xddqz4pbrbxbyi06sfl32ki0by90fxadcy";
  src = ./.;
}
