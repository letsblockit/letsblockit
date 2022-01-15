{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "sha256-MZKvii96fbqbxX+sPu8LYqt82/kojQJ9C2JNrnxPEq8=";
  src = ./.;
  doCheck = false;
}
