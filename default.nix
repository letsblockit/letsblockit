{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "0xy7vwabslsxrzairg23vw9iyqsv80bi563354gjrp2rkaj84xq5";
  src = ./.;
}
