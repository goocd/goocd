package targets

import "log"

func init() {

	TargetMap["atsame51-cmsisdap"] = TargetFunc(func(args *Args) error {

		log.Printf("GOT HERE: atsame51-cmsisdap")
		// Open File, Buffer here etc.
		err := args.Debugger.Program()
		if err != nil {
			return err
		}
		return nil
	})

}
