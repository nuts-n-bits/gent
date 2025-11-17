package main

func cgProgramTypescript(chkProgram LnkProgram, indent string) string {
	var b strings.Builder;
	b.WriteString(cgTsLib);
	for _, chkCommand := range chkProgram.commands {
		fragment := cgCommandTypescript(chkCommand, indent);
		b.WriteString(fragment);
	}
	for _, chkMix := range chkProgram.mixes {
		fragment := cgMixTypescript(chkMix, indent);
		b.WriteString(fragment);
	}
	return b.String();
}