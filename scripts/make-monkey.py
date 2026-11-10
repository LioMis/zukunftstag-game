#!/usr/bin/env python3
"""Builds assets/custom/monkey.glb: a cartoon monkey made of squashed spheres.

Run: python3 scripts/make-monkey.py   (needs numpy)

Units are metres: about 1.75 tall, feet at y=0, facing +Z like the trucks. One material per
colour (baseColorFactor, no texture). Edit PARTS and run again to change the look.
"""
import json
import math
import struct
from pathlib import Path

import numpy as np

FUR = (120, 72, 38)
SKIN = (232, 190, 140)
EYE = (25, 20, 20)
PATCH = (150, 95, 55)

# (centre, radii, colour) per sphere, in the monkey's own frame.
PARTS = [
    ((0, 0.75, 0), (0.45, 0.55, 0.4), FUR),  # body
    ((0, 0.72, 0.25), (0.3, 0.38, 0.18), SKIN),  # belly
    ((0, 1.45, 0), (0.38, 0.36, 0.36), FUR),  # head
    ((0, 1.48, -0.2), (0.26, 0.24, 0.18), PATCH),  # lighter fur on the back of the head
    ((0, 1.4, 0.24), (0.28, 0.24, 0.16), SKIN),  # face
    ((0, 1.32, 0.33), (0.17, 0.11, 0.1), SKIN),  # muzzle
    ((-0.11, 1.5, 0.36), (0.06, 0.07, 0.04), EYE),  # eyes
    ((0.11, 1.5, 0.36), (0.06, 0.07, 0.04), EYE),
    ((-0.38, 1.5, 0), (0.15, 0.15, 0.07), FUR),  # ears
    ((0.38, 1.5, 0), (0.15, 0.15, 0.07), FUR),
    ((-0.41, 1.5, 0.03), (0.09, 0.09, 0.05), SKIN),  # inner ears
    ((0.41, 1.5, 0.03), (0.09, 0.09, 0.05), SKIN),
    ((-0.44, 0.8, 0), (0.12, 0.38, 0.12), FUR),  # arms
    ((0.44, 0.8, 0), (0.12, 0.38, 0.12), FUR),
    ((-0.2, 0.28, 0), (0.14, 0.3, 0.14), FUR),  # legs
    ((0.2, 0.28, 0), (0.14, 0.3, 0.14), FUR),
]

# Tail: overlapping small spheres that leave the bottom, sweep back and up and curl in at the tip
# like a question mark. Thick at the base, thin at the tip.
N_TAIL = 40
for i in range(N_TAIL):
    t = i / (N_TAIL - 1)
    a = 1.03 - t * 4.3
    r = 0.58 * (1 - 0.6 * t)
    w = 0.085 - 0.04 * t
    PARTS.append(((0.1 * t, 0.8 - r * math.cos(a), -0.9 + r * math.sin(a)), (w, w, w), FUR))


def ellipsoid(centre, radii, rings=12, slices=16):
    """Positions, normals and triangle indices of a UV sphere scaled to radii."""
    pos, nrm = [], []
    for i in range(rings + 1):
        th = math.pi * i / rings
        for j in range(slices + 1):
            ph = 2 * math.pi * j / slices
            u = np.array([math.sin(th) * math.cos(ph), math.cos(th), math.sin(th) * math.sin(ph)])
            pos.append(np.array(centre) + u * radii)
            n = u / np.array(radii)
            nrm.append(n / np.linalg.norm(n))
    idx = []
    for i in range(rings):
        for j in range(slices):
            a = i * (slices + 1) + j
            b = a + slices + 1
            idx += [a, a + 1, b, a + 1, b + 1, b]
    return np.array(pos, np.float32), np.array(nrm, np.float32), np.array(idx, np.uint32)


def main():
    by_col = {}
    for centre, radii, col in PARTS:
        small = max(radii) < 0.1  # tail beads and eyes need fewer triangles
        by_col.setdefault(col, []).append(ellipsoid(centre, radii, *((6, 10) if small else (12, 16))))

    blob = bytearray()
    views, accessors, prims, mats = [], [], [], []

    def add(data, target, comp, typ, minmax=False):
        while len(blob) % 4:
            blob.append(0)
        views.append({"buffer": 0, "byteOffset": len(blob), "byteLength": data.nbytes, "target": target})
        blob.extend(data.tobytes())
        acc = {"bufferView": len(views) - 1, "componentType": comp, "count": len(data), "type": typ}
        if minmax:
            acc["min"] = data.min(axis=0).tolist()
            acc["max"] = data.max(axis=0).tolist()
        accessors.append(acc)
        return len(accessors) - 1

    for col, parts in by_col.items():
        pos, nrm, idx, off = [], [], [], 0
        for p, n, i in parts:
            pos.append(p)
            nrm.append(n)
            idx.append(i + off)
            off += len(p)
        pos, nrm, idx = np.concatenate(pos), np.concatenate(nrm), np.concatenate(idx)
        mats.append({"pbrMetallicRoughness": {"baseColorFactor": [c / 255 for c in col] + [1],
                                              "metallicFactor": 0, "roughnessFactor": 1}})
        prims.append({
            "attributes": {"POSITION": add(pos, 34962, 5126, "VEC3", True),
                           "NORMAL": add(nrm, 34962, 5126, "VEC3")},
            "indices": add(idx, 34963, 5125, "SCALAR"),
            "material": len(mats) - 1,
        })
    while len(blob) % 4:
        blob.append(0)

    gltf = {
        "asset": {"version": "2.0", "generator": "scripts/make-monkey.py"},
        "scene": 0, "scenes": [{"nodes": [0]}], "nodes": [{"mesh": 0, "name": "monkey"}],
        "meshes": [{"name": "monkey", "primitives": prims}],
        "materials": mats, "accessors": accessors, "bufferViews": views,
        "buffers": [{"byteLength": len(blob)}],
    }
    js = json.dumps(gltf, separators=(",", ":")).encode()
    js += b" " * (-len(js) % 4)
    out = Path(__file__).resolve().parent.parent / "assets" / "custom" / "monkey.glb"
    out.parent.mkdir(parents=True, exist_ok=True)
    with open(out, "wb") as f:
        f.write(struct.pack("<III", 0x46546C67, 2, 12 + 8 + len(js) + 8 + len(blob)))
        f.write(struct.pack("<II", len(js), 0x4E4F534A) + js)
        f.write(struct.pack("<II", len(blob), 0x004E4942) + bytes(blob))
    print("wrote", out, f"({out.stat().st_size // 1024} KB)")


if __name__ == "__main__":
    main()
