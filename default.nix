{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "1fwcaly2x79zw9wpm5v7xxd6q5k8kgbx6drncn84hqcxmvww7kqv";
  src = ./.;
}
