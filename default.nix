{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "0igjj37zcnh1gzmslllchyh7gybp1xgakbyfhj7vaqcvhnzl9822";
  src = ./.;
}
