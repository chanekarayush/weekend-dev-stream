import javax.swing.*;
import java.awt.*;

public class Demo {

    public static void main(String[] args) {
        JFrame mainFrame = new JFrame("The Person who lives is the person who lives");
        Container pane = new JFrame();
        pane.setLayout(new BorderLayout());

        // Controls Rotation
        JSlider horizontalSlider = new JSlider(0, 360, 180);
        pane.add(horizontalSlider, BorderLayout.SOUTH);
        JSlider verticalSlider = new JSlider(SwingConstants.VERTICAL, -180, 180, 0);
        pane.add(verticalSlider, BorderLayout.EAST);

        JPanel renderPanel = new JPanel() {
            public void paintComponents(Graphics g) {
            }
        };

        mainFrame.setSize(400, 400);
        mainFrame.setVisible(true);
    }

}
