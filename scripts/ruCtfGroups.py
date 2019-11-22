with open("teams","w") as f:
    for i in range(298):
        A=int(60+i/256)
        
        B=i%256
        f.write("10.{}.{}.2\n".format(str(A),str(B)))

    
