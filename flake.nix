{
  description = "Empiryczne wskaźniki niezawodności - dev shell (Go + Fyne GUI)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Fyne uses cgo + OpenGL + X11/Wayland, so we need the usual
        # native dev libraries on top of the Go toolchain.
        runtimeLibs = with pkgs; [
          libGL
          libx11
          libxcursor
          libxrandr
          libxinerama
          libxi
          libxext
          libxxf86vm
          wayland
          libxkbcommon
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            pkg-config
          ] ++ runtimeLibs;

          # Make sure cgo/gcc and dynamic libs are found at build and run time.
          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath runtimeLibs;
          CGO_ENABLED = "1";
        };
      });
}
