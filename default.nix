{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "0dv610gilwwqlchqhhxy74717s16v8dzirp6bisvn5ga1qh96fhw";
  src = ./.;
}
