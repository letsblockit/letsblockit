{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override{
  go = pkgs.go_1_16;
} {
  pname = "weblock";
  version = "1.0";
  vendorSha256 = "0qinf1gzfvv1mid6ayjg7dl10li0zppfm000knyrfixpl7s6xc5r";
  src = ./.;
}
