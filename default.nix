{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "1k8mp77ymm9k6i0bbfhsllxlqkz7zls83wg1kixslg2z1c27v40l";
  src = ./.;
}
