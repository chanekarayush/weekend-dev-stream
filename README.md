# weekend-dev-stream

## Week 3

### 3D Renderer from Scratch in Java Swing

A Java 3D Renderer that renders triangles to create a tetrahedron and a circle.
Also includes good stuff like rotation around both axes, Z Buffer Clipping and Shading.

### Matrix Representation

We use the `Matrix3d` class to store a 3×3 transformation matrix as a flat array of 9 values:

$$
\begin{bmatrix}
m_{00} & m_{01} & m_{02} \\
m_{10} & m_{11} & m_{12} \\
m_{20} & m_{21} & m_{22}
\end{bmatrix}
$$

---

## Heading Rotation (Yaw / Tranformation in X-Z Plane)

The horizontal slider controls rotation around the **X-Z Plane** also called as Yaw.

### Formula: Standard Y-Rotation Matrix

$$
R_y(\theta) =
\begin{bmatrix}
\cos(\theta) & 0 & -\sin(\theta) \\
0 & 1 & 0 \\
\sin(\theta) & 0 & \cos(\theta)
\end{bmatrix}
$$

### Effect on Vertex (x, y, z):

$$
\begin{bmatrix}
x' \\
y' \\
z'
\end{bmatrix} =
R_y(\theta) \cdot
\begin{bmatrix}
x \\
y \\
z
\end{bmatrix}
$$

Which expands to:

- $x' = x\cos(\theta) - z\sin(\theta)$
- $y' = y$ (unchanged, rotation around Y-axis)
- $z' = x\sin(\theta) + z\cos(\theta)$

---

## Pitch Rotation (Transformation along the Y-Z Plane)

The vertical slider controls tilt/rotation around the **Y-Z Plane** (pitch).

### Formula: Standard Y-Z Transformation Matrix

$$
R_x(\phi) =
\begin{bmatrix}
1 & 0 & 0 \\
0 & \cos(\phi) & -\sin(\phi) \\
0 & \sin(\phi) & \cos(\phi)
\end{bmatrix}
$$

### Effect on Vertex (x, y, z):

$$
\begin{bmatrix}
x'' \\
y'' \\
z''
\end{bmatrix} =
R_x(\phi) \cdot
\begin{bmatrix}
x' \\
y' \\
z'
\end{bmatrix}
$$

Which expands to:

- $x'' = x'$ (unchanged, rotation around X-axis)
- $y'' = y'\cos(\phi) - z'\sin(\phi)$
- $z'' = y'\sin(\phi) + z'\cos(\phi)$

---

## Combined Transformation

We take advantage of the **Associative Property** of Matrices where
$A.(B.C) = (A.B).C$
The final transformation matrix is the **product** of both matrices:

$$
T_{total} = R_y(\theta) \cdot R_x(\phi)
$$

### Matrix Multiplication (Code Implementation):

This applies rotations in sequence: first Y-axis rotation, then X-axis rotation. The resulting matrix $T_{total}$ can be computed as:

$$
T_{total} =
\begin{bmatrix}
\cos(\theta) & 0 & -\sin(\theta) \\
0 & 1 & 0 \\
\sin(\theta) & 0 & \cos(\theta)
\end{bmatrix} \cdot
\begin{bmatrix}
1 & 0 & 0 \\
0 & \cos(\phi) & -\sin(\phi) \\
0 & \sin(\phi) & \cos(\phi)
\end{bmatrix}
$$

### Final Matrix Elements:

- $m_{00} = \cos(\theta)$
- $m_{01} = 0$
- $m_{02} = -\sin(\theta)$
- $m_{10} = 0$
- $m_{11} = \cos(\phi)$
- $m_{12} = -\sin(\phi)$
- $m_{20} = \sin(\theta)$
- $m_{21} = \sin(\theta)\sin(\phi) + \cos(\theta)\cos(\phi)$ ← Combined effect!
- $m_{22} = -\sin(\theta)\sin(\phi) + \cos(\theta)\cos(\phi)$

---

## Transformation Applied to Vertex

For each triangle vertex (x, y, z):

$$
\begin{bmatrix}
x_{final} \\
y_{final} \\
z_{final}
\end{bmatrix} =
T_{total} \cdot
\begin{bmatrix}
x \\
y \\
z
\end{bmatrix}
$$

In code (Matrix3d.transform method):

$$
\begin{bmatrix}
x_{final} \\
y_{final} \\
z_{final}
\end{bmatrix} =
\begin{bmatrix}
m_{00} & m_{01} & m_{02} \\
m_{10} & m_{11} & m_{12} \\
m_{20} & m_{21} & m_{22}
\end{bmatrix} \cdot
\begin{bmatrix}
x \\
y \\
z
\end{bmatrix}
$$

---

## Rotation Matrix Properties

| Property          | Value/Formula        | Explanation                                  |
| ----------------- | -------------------- | -------------------------------------------- |
| **Determinant**   | +1 for both matrices | Pure rotation (no reflection)                |
| **Orthogonality** | $R^T R = I$          | Columns are orthonormal unit vectors         |
| **Inverse**       | $R^{-1} = R^T$       | Rotation matrix is its own transpose inverse |

---

### References

1. [**Rotation Matrix**](https://en.wikipedia.org/wiki/Rotation_matrix) - Standard linear algebra for 2D/3D rotations
2. [**Matrix Multiplication Order**](https://en.wikipedia.org/wiki/Multiplicative_group_of_rotation_matrices#Rotation_order) - $R_y(\theta) \cdot R_x(\phi)$ vs $R_x(\phi) \cdot R_y(\theta)$ produce different results
