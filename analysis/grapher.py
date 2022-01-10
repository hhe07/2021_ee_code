import csv
import matplotlib.pyplot as plt
import numpy as np
import scipy.stats
import sympy as sym
import matplotlib.ticker as ticker

name = "concat2"
from matplotlib import rcParams
rcParams['font.family'] = 'Liberation Serif'
rcParams['text.usetex'] = True

def scifmt(x):
    ret = "{:.4e}".format(x)
    ret = ret.split("e")
    if ret[1][0] == "-":
        ret[1] = ret[1].replace("0", "")
    return ret[0] + '*10^{' + ret[1] + '}'

def logfit(x, y):
    slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(np.log(x), y)
    
    label = '$y={} \ln x + {}, R^2 = {:.4f}$'.format(scifmt(slope), scifmt(intercept), r_value)
    return slope * np.log(x) + intercept, label

def linfit(x,y):
    slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(x, y)
    label = '$y={}x + {}, R^2 = {:.4f}$'.format(scifmt(slope), scifmt(intercept), r_value)
    return slope * x + intercept, label


with open(name+".csv", newline = '') as csvfile:
    res = csv.reader(csvfile)
    names = []
    scale = []
    data = [[],[],[],[]]
    #      rope, rope err, gap, gap err
    ct = 0
    for row in res:
        if ct == 0:
            names = row
            ct+=1
        else:
            idx = 0
            scale.append(float(row[0]))
            for el in row[1:]:
                data[idx].append(float(el))
                idx +=1
    fig = plt.figure()
    dep = fig.add_subplot(111)
    dep.errorbar(scale, data[0], data[1], fmt = "o", label = "Rope", linewidth = 2, capsize = 4)
    dep.errorbar(scale, data[2], data[3], fmt = "o", label = "Gap Buffer",linewidth = 2,  capsize = 4)
    resr, lr = linfit(np.asarray(scale), data[0])
    plt.plot(scale, resr,label = "Rope Regression: "+ lr)
    
    resg, lg = linfit(np.asarray(scale), data[2])
    plt.plot(scale, resg,label = "Gap Buffer Regression: "+ lg)
    
    lg = plt.legend(bbox_to_anchor = (0.5, -0.3), loc = "center")

    
    plt.suptitle(names[0])
    plt.xlabel(names[1])
    plt.ylabel("Avg. Time / Operation (ns)")
    if len(names) > 2:
        plt.title(names[2])
    
    plt.savefig(name + '.png', 
                dpi=300, 
                format='png', 
                bbox_extra_artists=(lg,), 
                bbox_inches='tight',
                pad_inches = 0.5)
    
    #plt.show()
    
    
    