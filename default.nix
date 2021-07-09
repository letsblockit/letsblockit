{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "07lk7pbfvyf3y3nk097jq494kcq7cm5al6kpigcgfyzqa3rp5m31";
  src = ./.;
}
