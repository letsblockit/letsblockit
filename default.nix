{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "16rfbn96bfyxbmrnh8qj02qcf9ggkc72zgd6gs62zp1271pm4xnx";
  src = ./.;
}
