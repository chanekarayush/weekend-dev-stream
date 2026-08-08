# Weekend Dev Stream

## Week 2

### Video Codec in Golang


This program is a toy codec that is based on real video codecs, just a lot less math.

### Overview of Codec Process

<img width="1884" height="618" alt="image" src="https://github.com/user-attachments/assets/bc1c197a-9d1a-4c98-b2a6-308bc7db5e82" />

The Above diagram roughly outlines the entire process of conversion and encoding and decoding.


### Chroma Subsampling

What I mean by this is not just Subsampling which we do here using the most popular 4:2:0 
ratio, but the entire process of converting RGB to YUV colorspace
<img width="1444" height="659" alt="image" src="https://github.com/user-attachments/assets/d4f5a21c-a136-4921-887c-46e167248cfc" />

#### RGB to YUV conversion formulae
All the math in this program is basically reduced to these 3 formulae
<img width="699" height="179" alt="image" src="https://github.com/user-attachments/assets/d167182b-ca44-4808-ad6e-5edf11861418" />
*Source: [Wikipedia](https://en.wikipedia.org/wiki/Y%E2%80%B2UV )*

### References

1. [kevmo314's Video Codec Implementation](https://github.com/kevmo314/codec-from-scratch/)
2. [Y' UV Colorspace](https://en.wikipedia.org/wiki/Y%E2%80%B2UV )
3. [Chroma Subsampling](https://en.wikipedia.org/wiki/Chroma_subsampling#Gamut_clipping)
4. [Key frames](https://en.wikipedia.org/wiki/Key_frame )
5. [Run-length Encoding](https://en.wikipedia.org/wiki/Run-length_encoding )
