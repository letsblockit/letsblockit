{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0xchc8f6adhdsi2ncf7h0lr43fa2mvzddkhf9lbc49mlgnnvxhgx";
  src = ./.;
}
