// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Myedgeorder.Outputs
{

    /// <summary>
    /// Different types of filters supported and its values.
    /// </summary>
    [OutputType]
    public sealed class FilterablePropertyResponse
    {
        /// <summary>
        /// Values to be filtered.
        /// </summary>
        public readonly ImmutableArray<string> SupportedValues;
        /// <summary>
        /// Type of product filter.
        /// </summary>
        public readonly string Type;

        [OutputConstructor]
        private FilterablePropertyResponse(
            ImmutableArray<string> supportedValues,

            string type)
        {
            SupportedValues = supportedValues;
            Type = type;
        }
    }
}