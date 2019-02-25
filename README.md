# HexaXAnd)

## Youtube demo video
[![](http://img.youtube.com/vi/-3kaH5sj-EM/0.jpg)](http://www.youtube.com/watch?v=-3kaH5sj-EM "HexaXAnd0 Demo")

## Screens
![Demo Image ](https://github.com/msorins/HexaXAnd0/blob/GameDetection/0.jpeg?raw=true "Demo Image")

![Demo Image ](https://github.com/msorins/HexaXAnd0/blob/GameDetection/1.jpeg?raw=true "Demo Image")

![Demo Image ](https://github.com/msorins/HexaXAnd0/blob/GameDetection/2.jpeg?raw=true "Demo Image")

# Idea
Playing X and 0 with a robot.

Place the playing board somewhere on a wall (in the reaching distance of the robot), take turns and try to win the game

# How does it work
I am using OpenCV in order to process the image and get bounding boxes for the game board itself and also for each of the smaller 9 squares (where *X* or *0* would be written).

In more details, the image processing pipeline looks like this:

1. Inverting pixels 
2. Applying threshold
3. Computing contours
4. Filtering contours
    * Approximating the contours (to remove extra edges)
    * Keeping only the contours that have 4 edges
    * Looking for a bigger contour that includes 9 sub contours (game board with all the 9 smaller squares)
5. Iterating through contours and looking at the pixels colours to determine if a square is empty, contains and X or 0 (the robot and the player are going to use different colours)

# Robot writing
I have mounted a marker on the robot's head and using a built in distance sensor it will approach the wall until the tip of the marker touches it, and then it nod its head or tilt its body to make an *X*

*Problem:* the distance to the wall must be very precise, otherwise the marker would just bend

# Technologies used

* GO
* OpenCV


> The Project was realised in the 3rd year of University