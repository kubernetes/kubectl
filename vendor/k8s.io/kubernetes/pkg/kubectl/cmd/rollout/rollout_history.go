/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rollout

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions/printers"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions/resource"
	"k8s.io/kubernetes/pkg/kubectl/polymorphichelpers"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

var (
	history_long = templates.LongDesc(`
		View previous rollout revisions and configurations.`)

	history_example = templates.Examples(`
		# View the rollout history of a deployment
		kubectl rollout history deployment/abc

		# View the details of daemonset revision 3
		kubectl rollout history daemonset/abc --revision=3`)
)

type RolloutHistoryOptions struct {
	PrintFlags *genericclioptions.PrintFlags
	ToPrinter  func(string) (printers.ResourcePrinter, error)

	Revision int64

	Builder          func() *resource.Builder
	Resources        []string
	Namespace        string
	EnforceNamespace bool

	HistoryViewer    polymorphichelpers.HistoryViewerFunc
	RESTClientGetter genericclioptions.RESTClientGetter

	resource.FilenameOptions
	genericclioptions.IOStreams
}

func NewRolloutHistoryOptions(streams genericclioptions.IOStreams) *RolloutHistoryOptions {
	return &RolloutHistoryOptions{
		PrintFlags: genericclioptions.NewPrintFlags("").WithTypeSetter(scheme.Scheme),
		IOStreams:  streams,
	}
}

func NewCmdRolloutHistory(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewRolloutHistoryOptions(streams)

	validArgs := []string{"deployment", "daemonset", "statefulset"}

	cmd := &cobra.Command{
		Use: "history (TYPE NAME | TYPE/NAME) [flags]",
		DisableFlagsInUseLine: true,
		Short:   i18n.T("View rollout history"),
		Long:    history_long,
		Example: history_example,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run())
		},
		ValidArgs: validArgs,
	}

	cmd.Flags().Int64Var(&o.Revision, "revision", o.Revision, "See the details, including podTemplate of the revision specified")

	usage := "identifying the resource to get from a server."
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)

	o.PrintFlags.AddFlags(cmd)

	return cmd
}

func (o *RolloutHistoryOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.Resources = args

	var err error
	if o.Namespace, o.EnforceNamespace, err = f.ToRawKubeConfigLoader().Namespace(); err != nil {
		return err
	}

	o.ToPrinter = func(operation string) (printers.ResourcePrinter, error) {
		o.PrintFlags.NamePrintFlags.Operation = operation
		return o.PrintFlags.ToPrinter()
	}

	o.HistoryViewer = polymorphichelpers.HistoryViewerFn
	o.RESTClientGetter = f
	o.Builder = f.NewBuilder

	return nil
}

func (o *RolloutHistoryOptions) Validate() error {
	if len(o.Resources) == 0 && cmdutil.IsFilenameSliceEmpty(o.Filenames) {
		return fmt.Errorf("Required resource not specified.")
	}
	if o.Revision < 0 {
		return fmt.Errorf("revision must be a positive integer: %v", o.Revision)
	}

	return nil
}

func (o *RolloutHistoryOptions) Run() error {

	r := o.Builder().
		WithScheme(legacyscheme.Scheme).
		NamespaceParam(o.Namespace).DefaultNamespace().
		FilenameParam(o.EnforceNamespace, &o.FilenameOptions).
		ResourceTypeOrNameArgs(true, o.Resources...).
		ContinueOnError().
		Latest().
		Flatten().
		Do()
	if err := r.Err(); err != nil {
		return err
	}

	return r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		mapping := info.ResourceMapping()
		historyViewer, err := o.HistoryViewer(o.RESTClientGetter, mapping)
		if err != nil {
			return err
		}
		historyInfo, err := historyViewer.ViewHistory(info.Namespace, info.Name, o.Revision)
		if err != nil {
			return err
		}

		withRevision := ""
		if o.Revision > 0 {
			withRevision = fmt.Sprintf("with revision #%d", o.Revision)
		}

		printer, err := o.ToPrinter(fmt.Sprintf("%s\n%s", withRevision, historyInfo))
		if err != nil {
			return err
		}

		return printer.PrintObj(cmdutil.AsDefaultVersionedOrOriginal(info.Object, info.Mapping), o.Out)
	})
}
