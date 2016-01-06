#!/usr/bin/env python
# coding: utf8

import os
import sys
import re

DECLARE_PATH_PATTER = re.compile(r'# (.*)')
ERROR_PATTERN = re.compile(r'(.*):(\d+): (.*)')

fdin, fdout, fderr = os.popen3('gb build')
try:
	sys.stdout.write(fdout.read()) # just output stdout
	errlines = fderr.readlines()
finally:
	fdin.close()
	fdout.close()
	fderr.close()

buildingPath = None
for line in errlines:
	line = line.strip('\n')

	m = DECLARE_PATH_PATTER.match(line)
	if m: 
		buildingPath = m.groups()[0]
		print >>sys.stderr, line
		continue 

	m = ERROR_PATTERN.match(line)
	if m: 
		fn, lno, errmsg = m.groups()
		print >>sys.stderr, '%s/%s:%s: %s' % (buildingPath, fn, lno, errmsg)
		continue 

	print >>sys.stderr, line
	continue 

