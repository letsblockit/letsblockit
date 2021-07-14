{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "01dz206wbnxss3vylzs6mp3xa4ip6mymk3v8bxxhx178anjak7i1";
  src = ./.;
}
