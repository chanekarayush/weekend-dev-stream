import javax.swing.*;
import java.awt.*;
import java.awt.geom.Path2D;
import java.awt.image.BufferedImage;
import java.util.ArrayList;
import java.util.List;

public class Demo {

    public static void main(String[] args) {
        JFrame mainFrame = new JFrame();
        Container pane = mainFrame.getContentPane();
        pane.setLayout(new BorderLayout());

        // Controls Rotation
        JSlider horizontalSlider = new JSlider(0, 360, 180);
        pane.add(horizontalSlider, BorderLayout.SOUTH);

        JSlider verticalSlider = new JSlider(SwingConstants.VERTICAL, -90, 90, 0);
        pane.add(verticalSlider, BorderLayout.EAST);

        JPanel renderPanel = new JPanel() {
            public void paintComponent(Graphics g) {
                Graphics2D g2 = (Graphics2D) g;
                g2.setColor(Color.BLACK);
                g2.fillRect(0, 0, getWidth(), getHeight());

                List<Triangle> triangles = new ArrayList<>();

                triangles.add(
                        new Triangle(new Vertex(100, 100, 100),
                                new Vertex(-100, -100, 100),
                                new Vertex(-100, 100, -100),
                                Color.WHITE));
                triangles.add(
                        new Triangle(new Vertex(100, 100, 100),
                                new Vertex(-100, -100, 100),
                                new Vertex(100, -100, -100),
                                Color.MAGENTA));
                triangles.add(
                        new Triangle(new Vertex(-100, 100, -100),
                                new Vertex(100, -100, -100),
                                new Vertex(100, 100, 100),
                                Color.CYAN));
                triangles.add(
                        new Triangle(new Vertex(-100, 100, -100),
                                new Vertex(100, -100, -100),
                                new Vertex(-100, -100, 100),
                                Color.YELLOW));
                double heading = Math.toRadians(horizontalSlider.getValue());
                Matrix3d headingTransform = new Matrix3d(new double[] {
                        Math.cos(heading), 0, -Math.sin(heading),
                        0, 1, 0,
                        Math.sin(heading), 0, Math.cos(heading)
                });

                double pitch = Math.toRadians(verticalSlider.getValue());
                Matrix3d pitchTransform = new Matrix3d(new double[] {
                        1, 0, 0,
                        0, Math.cos(pitch), Math.sin(pitch),
                        0, -Math.sin(pitch), Math.cos(pitch)
                });

                Matrix3d transform = headingTransform.multiply(pitchTransform);

                BufferedImage img = new BufferedImage(getWidth(), getHeight(), BufferedImage.TYPE_INT_ARGB);
                double[] zBuffer = new double[getWidth() * getHeight()];

                for (int i = 0; i < zBuffer.length; i++) {
                    zBuffer[i] = Double.NEGATIVE_INFINITY;
                }

                for (Triangle t : triangles) {
                    Vertex v1 = transform.transform(t.v1);
                    Vertex v2 = transform.transform(t.v2);
                    Vertex v3 = transform.transform(t.v3);

                    v1.x += getWidth() / 2;
                    v2.x += getWidth() / 2;
                    v3.x += getWidth() / 2;

                    v1.y += getHeight() / 2;
                    v2.y += getHeight() / 2;
                    v3.y += getHeight() / 2;

                    Vertex ab = new Vertex(v2.x - v1.x, v2.y - v1.y, v2.z - v1.z);
                    Vertex ac = new Vertex(v3.x - v1.x, v3.y - v1.y, v3.z - v1.z);

                    Vertex norm = new Vertex(
                            ab.y * ac.z - ab.z * ac.y,
                            ab.z * ac.x - ab.x * ac.z,
                            ab.x * ac.y - ab.y * ac.x);

                    double normalLength = Math.sqrt(norm.x * norm.x + norm.y * norm.y + norm.z * norm.z);
                    norm.x /= normalLength;
                    norm.y /= normalLength;
                    norm.z /= normalLength;

                    double angleCos = Math.abs(norm.z);

                    int minX = (int) Math.max(0, Math.min(v1.x, Math.min(v2.x, v3.x)));
                    int minY = (int) Math.max(0, Math.min(v1.y, Math.min(v2.y, v3.y)));
                    int maxX = (int) Math.min(img.getWidth() - 1, Math.max(v1.x, Math.max(v2.x, v3.x)));
                    int maxY = (int) Math.min(img.getHeight() - 1, Math.max(v1.y, Math.max(v2.y, v3.y)));

                    double triangleArea = (v1.y - v3.y) * (v2.x - v3.x) + (v2.y - v3.y) * (v3.x - v1.x);

                    for (int y = minY; y <= maxY; y++) {
                        for (int x = minX; x <= maxX; x++) {
                            double b1 = ((y - v3.y) * (v2.x - v3.x) + (v2.y - v3.y) * (v3.x - x)) / triangleArea;
                            double b2 = ((y - v1.y) * (v3.x - v1.x) + (v3.y - v1.y) * (v1.x - x)) / triangleArea;
                            double b3 = ((y - v2.y) * (v1.x - v2.x) + (v1.y - v2.y) * (v2.x - x)) / triangleArea;
                            if (b1 >= 0 && b1 <= 1 && b2 >= 0 && b2 <= 1 && b3 >= 0 && b3 <= 1) {
                                int zIndex = y * img.getHeight() + x;
                                double depth = b1 * v1.z + b2 * v2.z + b3 * v3.z;
                                if (zBuffer[zIndex] < depth) {
                                    img.setRGB(x, y, getShade(t.clr, angleCos).getRGB());
                                    zBuffer[zIndex] = depth;
                                }
                            }
                        }
                    }
                }
                g2.drawImage(img, 0, 0, null);
            }
        };
        horizontalSlider.addChangeListener(e -> renderPanel.repaint());
        verticalSlider.addChangeListener(e -> renderPanel.repaint());

        pane.add(renderPanel, BorderLayout.CENTER);

        mainFrame.setSize(400, 400);
        mainFrame.setVisible(true);
        mainFrame.setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
        mainFrame.setTitle("Triangle");
    }

    public static Color getShade(Color color, double shade) {
        double redLinear = Math.pow(color.getRed(), 2.4) * shade;
        double greenLinear = Math.pow(color.getGreen(), 2.4) * shade;
        double blueLinear = Math.pow(color.getBlue(), 2.4) * shade;

        int red = (int) Math.pow(redLinear, 1 / 2.4);
        int green = (int) Math.pow(greenLinear, 1 / 2.4);
        int blue = (int) Math.pow(blueLinear, 1 / 2.4);

        return new Color(red, green, blue);
    }

}

/**
 * Triangle
 * Holds a triangle made from Vertex
 */
class Triangle {
    Vertex v1;
    Vertex v2;
    Vertex v3;
    Color clr;

    Triangle(Vertex v1, Vertex v2, Vertex v3, Color clr) {
        this.v1 = v1;
        this.v2 = v2;
        this.v3 = v3;
        this.clr = clr;
    }
}

/**
 * Vertex
 */
class Vertex {
    double x;
    double y;
    double z;

    Vertex(double x, double y, double z) {
        this.x = x;
        this.y = y;
        this.z = z;
    }
}

/**
 * Matrix3d
 */
class Matrix3d {
    double[] values;

    Matrix3d(double[] values) {
        this.values = values;
    }

    Matrix3d multiply(Matrix3d other) {
        double[] result = new double[9];
        for (int r = 0; r < 3; r++) {
            for (int c = 0; c < 3; c++) {
                for (int i = 0; i < 3; i++) {
                    result[r * 3 + c] += this.values[r * 3 + i] * other.values[i * 3 + c];
                }
            }
        }
        return new Matrix3d(result);
    }

    Vertex transform(Vertex v) {
        return new Vertex(
                v.x * values[0] + v.y * values[1] + v.z * values[2],
                v.x * values[3] + v.y * values[4] + v.z * values[5],
                v.x * values[6] + v.y * values[7] + v.z * values[8]);
    }
}
